package routes

import (
	"database/sql"
	"net/http"

	"github.com/kwizera/university-invoices/internal/middleware"
	feecontrollers "github.com/kwizera/university-invoices/services/fee/controllers"
	invoicecontrollers "github.com/kwizera/university-invoices/services/invoice/controllers"
	paymentcontrollers "github.com/kwizera/university-invoices/services/payment/controllers"
	studentcontrollers "github.com/kwizera/university-invoices/services/student/controllers"
	usercontrollers "github.com/kwizera/university-invoices/services/user/controllers"
	"github.com/kwizera/university-invoices/shared/repositories"
)

func SetupRouter(db *sql.DB) http.Handler {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	studentRepo := repositories.NewStudentRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)
	invoiceItemRepo := repositories.NewInvoiceItemRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	feeRepo := repositories.NewFeeRepository(db)
	receiptRepo := repositories.NewReceiptRepository(db)

	// Initialize controllers
	userController := usercontrollers.NewUserController(userRepo)
	studentController := studentcontrollers.NewStudentController(studentRepo)
	invoiceController := invoicecontrollers.NewInvoiceController(invoiceRepo, invoiceItemRepo, studentRepo)
	paymentController := paymentcontrollers.NewPaymentController(paymentRepo, invoiceRepo, receiptRepo)
	feeController := feecontrollers.NewFeeController(feeRepo)

	// Create router
	mux := http.NewServeMux()

	// Apply CORS middleware
	handler := middleware.CORS(mux)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Auth routes
	mux.HandleFunc("/api/auth/login", userController.Login)
	mux.HandleFunc("/api/auth/register", userController.CreateUser)

	// User routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userController.ListUsers(w, r)
		case http.MethodPost:
			userController.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userController.GetUser(w, r)
		case http.MethodPut:
			userController.UpdateUser(w, r)
		case http.MethodDelete:
			userController.DeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Student routes
	mux.HandleFunc("/api/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			studentController.ListStudents(w, r)
		case http.MethodPost:
			studentController.CreateStudent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/students/by-id", studentController.GetStudentByStudentID)
	mux.HandleFunc("/api/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			studentController.GetStudent(w, r)
		case http.MethodPut:
			studentController.UpdateStudent(w, r)
		case http.MethodDelete:
			studentController.DeleteStudent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Fee type routes
	mux.HandleFunc("/api/fees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			feeController.ListFeeTypes(w, r)
		case http.MethodPost:
			feeController.CreateFeeType(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/fees/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			feeController.GetFeeType(w, r)
		case http.MethodPut:
			feeController.UpdateFeeType(w, r)
		case http.MethodDelete:
			feeController.DeleteFeeType(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Invoice routes
	mux.HandleFunc("/api/invoices", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			invoiceController.ListInvoices(w, r)
		case http.MethodPost:
			invoiceController.CreateInvoice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/invoices/student", invoiceController.GetInvoicesByStudent)
	mux.HandleFunc("/api/invoices/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			invoiceController.GetInvoice(w, r)
		case http.MethodPut:
			invoiceController.UpdateInvoice(w, r)
		case http.MethodDelete:
			invoiceController.DeleteInvoice(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Payment routes
	mux.HandleFunc("/api/payments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			paymentController.ListPayments(w, r)
		case http.MethodPost:
			paymentController.CreatePayment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/payments/invoice", paymentController.GetPaymentsByInvoice)
	mux.HandleFunc("/api/payments/student", paymentController.GetPaymentsByStudent)
	mux.HandleFunc("/api/payments/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			paymentController.GetPayment(w, r)
		case http.MethodPut:
			paymentController.UpdatePayment(w, r)
		case http.MethodDelete:
			paymentController.DeletePayment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return handler
}
