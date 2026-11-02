-- Seed student profiles for the default student users
-- These correspond to the users created in migration 009
-- Emails and names are encrypted using AES-256-GCM with ENCRYPTION_KEY from environment
-- Default encryption key for development: 12345678901234567890123456789012

INSERT INTO students (student_id, first_name, last_name, email, phone, program, semester, year, address, is_active) VALUES
('STU001', 'PyVfXw0bdcapZFaThsv7kDRo1osPEmI5gNvH1QzutsHsPeKTdjuEk6WAWg==', '', 'TAwyve/0yfOyvp69d/OQ68ugRXbGQVF3tU5KIXSP6+glvjfWezufNAFemvf40koZGMJSntU2GAjEKBzt', '+250788123456', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE),
('STU002', '11mJogmd6fJjLjVNzM4bOIRHGLeEwAZrRN9s6EnvkRQNSuqe', 'lPLe2VaLPk8you1Y654KM0S7a4zKKTGjvJjjIHSf7GWCx2VyCIk=', 'yrAADhR5LXgO1GigARLKH7wexjVJ9pu8x2enMl7xZytKxFfvrxHEo//99F1FTw4LlMhS1Hc=', '+250788234567', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE),
('STU003', 'Xkz86ZGnvbOWUNDF+5WeuR7wrkOKHY0w6Eigqh14bTAYSQzR3sA=', '2Yvd4/+z6EyjCLRP0sExK40/zrMRALnQiSF0qqCyANK5G1DtiU8922A=', '4litX1/vVP+rR7zySgh/kTgUWC9tr1EY0ltkoGj4SFD/0XdbFLjZtHbupTuIu+kfou9pZAWE', '+250788345678', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE),
('STU004', 'Jjngnw4pO8Qo3dInJv6XOrOrVlFzU2uNInnwD2xUyINIeHw=', 'vMZF2PDEdZw4JdLwZD/b/V68gfVz2husnCUpSSGLm1c=', 'CLUXIzxC+SI1tdvFzN5uB1YTGTDHUaUnZsJo6zXnVM3XkHdeWXp0LAu9TE30GvA+pg==', '+250788456789', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE),
('STU005', 'OoJnvLpSZFUb3tNB3DBs0EjpQRyO4dBPJLGr03e3QVkQZgV9', 'OCIlMg05eaLLkJ7UvwgG10lSepEVzwSwor0oQiGCtVcugXwGNQ==', '3kJvWthICg0eZXp/MpQsiIxRJb1m2y+VTHvOHtI9nObVC6myE5DoVPfkbk/LS9h4TJ2COjw=', '+250788567890', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE),
('STU006', 'kdwgL+lC6/t94RiSrq6j7ZypZq2yXklbofHqqWghwD9G23aDsEU=', 'UmJFT/evK0+SWCkoQgaHrMGLCRnHUtY2btGbCmvqJZ7k9tui10Gltw==', 'Wvhd+vWFB/JgT5U6loKtaFBPXUFJg91LpSRoJtxcVnHAEbp/D3LvY3O4qpGK9wV6r1Y6', '+250788678901', 'Computer Science', 'Fall 2024', 2024, 'Kigali, Rwanda', TRUE);
