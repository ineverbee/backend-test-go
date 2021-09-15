package auto

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ineverbee/backend-test-go/internal/app/cache"
	"github.com/ineverbee/backend-test-go/internal/app/currency"
	"github.com/ineverbee/backend-test-go/internal/app/models"
	"github.com/ineverbee/backend-test-go/internal/app/store"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ctxKey int8

const (
	ctxKeyUser ctxKey = iota
	ctxKeyRequestID
)

// Server structure
type server struct {
	logger *logrus.Logger
	router *mux.Router
	store  store.Store
	cache  cache.Cache
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Store, cache cache.Cache) *server {
	s := &server{
		logger: logrus.New(),
		router: mux.NewRouter(),
		store:  store,
		cache:  cache,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.Use(s.accessLogMiddleware)
	s.router.Use(s.panicMiddleware)
	s.router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, map[string]bool{"ok": true})
	})
	s.router.HandleFunc("/api/balance", s.handleBalanceChange()).Methods("POST")
	s.router.HandleFunc("/api/transactions", s.handleTransactions()).Methods("GET")
	s.router.HandleFunc("/api", s.handleBalance()).Methods("GET")
}

func (s *server) handleBalanceChange() http.HandlerFunc {
	type request struct {
		ID        int64  `json:"id"`
		Amount    int    `json:"amount"`
		Operation string `json:"operation"`
		Person    *int64 `json:"person"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println(r.Body, req)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			ID: req.ID,
		}
		ut := &models.Transactions{
			BalanceChange: req.Amount,
		}
		p, pt := &models.User{}, &models.Transactions{}

		if req.Person != nil {
			p.ID = *req.Person
			pt.BalanceChange = req.Amount
			if req.Operation == "accrual" {
				ut.Comment = "transaction from id" + strconv.Itoa(int(*req.Person))
				pt.Comment = "transaction to id" + strconv.Itoa(int(req.ID))
			} else {
				ut.Comment = "transaction to id" + strconv.Itoa(int(*req.Person))
				pt.Comment = "transaction from id" + strconv.Itoa(int(req.ID))
			}
			if err := s.store.Users().FindByID(p); err == sql.ErrNoRows {
				s.error(w, r, http.StatusUnprocessableEntity, errors.New("non-existent person"))
				return
			} else if err != nil && err != sql.ErrNoRows {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			} else {
				if req.Operation == "write-off" {
					//Updates existing user
					p.Balance += req.Amount
					if err := s.store.Users().Update(p); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
				} else if req.Operation == "accrual" {
					if p.Balance >= req.Amount {
						p.Balance -= req.Amount
						if err := s.store.Users().Update(p); err != nil {
							s.error(w, r, http.StatusUnprocessableEntity, err)
							return
						}
					} else {
						s.error(w, r, http.StatusUnprocessableEntity, errors.New("insufficient funds"))
						return
					}
				} else {
					s.error(w, r, http.StatusUnprocessableEntity, errors.New("there is no such operation like "+req.Operation))
					return
				}
			}
		} else {
			ut.Comment = req.Operation
		}

		if err := s.store.Users().FindByID(u); err == sql.ErrNoRows {
			if req.Operation == "accrual" {
				// Creates user if doesn't exists one
				u.Balance = req.Amount
				if err := s.store.Users().Create(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			} else if req.Operation == "write-off" {
				s.error(w, r, http.StatusUnprocessableEntity, errors.New("cannot write off funds from a non-existent user"))
				return
			} else {
				s.error(w, r, http.StatusUnprocessableEntity, errors.New("there is no such operation like "+req.Operation))
				return
			}
		} else if err != nil && err != sql.ErrNoRows {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		} else {
			if req.Operation == "accrual" {
				//Updates existing user
				u.Balance += req.Amount
				if err := s.store.Users().Update(u); err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, err)
					return
				}
			} else if req.Operation == "write-off" {
				if u.Balance >= req.Amount {
					u.Balance -= req.Amount
					if err := s.store.Users().Update(u); err != nil {
						s.error(w, r, http.StatusUnprocessableEntity, err)
						return
					}
				} else {
					s.error(w, r, http.StatusUnprocessableEntity, errors.New("insufficient funds"))
					return
				}
			} else {
				s.error(w, r, http.StatusUnprocessableEntity, errors.New("there is no such operation like "+req.Operation))
				return
			}
		}
		ut.User_ID = u.ID
		pt.User_ID = p.ID
		s.store.Transactions().CreateTransaction(ut)
		s.store.Transactions().CreateTransaction(pt)
		s.respond(w, r, http.StatusOK, u)
		s.cache.Set(u.ID, u.Balance, time.Hour*24)
	}
}

func (s *server) handleBalance() http.HandlerFunc {
	type request struct {
		ID       int64  `json:"id"`
		Currency string `json:"currency"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &models.User{
			ID: req.ID,
		}
		if err := s.store.Users().FindByID(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if req.Currency != "" {
			f, err := currency.ChangeCurrency(req.Currency, u.Balance)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			s.respond(w, r, http.StatusOK, struct {
				ID        int64   `json:"id"`
				Balance   int     `json:"balance"`
				Currency  string  `json:"currency"`
				Converted float64 `json:"converted"`
			}{ID: u.ID, Balance: u.Balance, Currency: req.Currency, Converted: f})
			s.cache.Set(u.ID, u.Balance, time.Hour*24)
		} else {
			s.respond(w, r, http.StatusOK, u)
			s.cache.Set(u.ID, u.Balance, time.Hour*24)
		}
	}
}

func (s *server) handleTransactions() http.HandlerFunc {
	type request struct {
		ID int64 `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		t := &[]models.Transactions{}

		if err := s.store.Transactions().ListOfTransactions(req.ID, t); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		keys, ok := r.URL.Query()["sort"]

		if ok || len(keys) > 0 {
			sorting := keys[0]

			if sorting != "" {
				if sorting == "amount" {
					sort.Slice(*t, func(i, j int) bool {
						return (*t)[i].BalanceChange < (*t)[j].BalanceChange
					})
				} else if sorting == "date" {
					sort.Slice(*t, func(i, j int) bool {
						return (*t)[i].CreatedAt < (*t)[j].CreatedAt
					})
				} else {
					s.error(w, r, http.StatusUnprocessableEntity, errors.New("there is no such sort like:"+sorting))
					return
				}
			}
		}

		keys, ok = r.URL.Query()["page"]

		if ok || len(keys) > 0 {
			page := keys[0]

			if page != "" {
				i, err := strconv.Atoi(page)
				if err != nil {
					s.error(w, r, http.StatusUnprocessableEntity, errors.New("there is no such page like:"+page))
					return
				}
				paginate(t, (i-1)*10, 10*i)
			}
		}

		s.respond(w, r, http.StatusOK, t)
	}
}

func paginate(x *[]models.Transactions, skip int, size int) {
	if skip > len(*x) {
		skip = len(*x)
	}

	end := skip + size
	if end > len(*x) {
		end = len(*x)
	}

	*x = (*x)[skip:end]
}
