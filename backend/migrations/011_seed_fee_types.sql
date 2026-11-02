-- Seed default fee types
-- These are common fee types for University of Kigali

INSERT INTO fee_types (name, description, amount, currency, is_recurring, is_active) VALUES
('Tuition Fee', 'Semester tuition fee for Computer Science program', 500000.00, 'RWF', TRUE, TRUE),
('Registration Fee', 'Annual registration fee', 50000.00, 'RWF', TRUE, TRUE),
('Library Fee', 'Library access and resource fee', 25000.00, 'RWF', TRUE, TRUE),
('Laboratory Fee', 'Computer laboratory usage fee', 30000.00, 'RWF', TRUE, TRUE),
('Examination Fee', 'Semester examination fee', 20000.00, 'RWF', TRUE, TRUE),
('Student ID Card', 'One-time student ID card fee', 5000.00, 'RWF', FALSE, TRUE),
('Medical Insurance', 'Annual medical insurance fee', 15000.00, 'RWF', TRUE, TRUE),
('Sports Fee', 'Sports facilities and activities fee', 10000.00, 'RWF', TRUE, TRUE);
