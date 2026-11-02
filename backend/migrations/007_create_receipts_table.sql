-- Create receipts table
CREATE TABLE IF NOT EXISTS receipts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    receipt_number VARCHAR(50) NOT NULL UNIQUE,
    payment_id INT NOT NULL,
    invoice_id INT NOT NULL,
    student_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    issued_date DATE NOT NULL,
    issued_by INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE RESTRICT,
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE RESTRICT,
    FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE RESTRICT,
    FOREIGN KEY (issued_by) REFERENCES users(id) ON DELETE RESTRICT,
    INDEX idx_receipt_number (receipt_number),
    INDEX idx_payment_id (payment_id),
    INDEX idx_student_id (student_id),
    INDEX idx_issued_date (issued_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
