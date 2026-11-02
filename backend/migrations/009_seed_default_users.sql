-- Seed default student users
-- Passwords are hashed using bcrypt (default: password123)
-- Emails and names are encrypted using AES-256-GCM with ENCRYPTION_KEY from environment
-- Default password for all students: password123
-- Bcrypt hash for 'password123': $2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C
-- Default encryption key for development: 12345678901234567890123456789012
-- IMPORTANT: Set ENCRYPTION_KEY in .env file for production (must be exactly 32 characters)

-- Student users (emails and names are encrypted)
INSERT INTO users (email, password, first_name, last_name, role, student_id, is_active) VALUES
('TAwyve/0yfOyvp69d/OQ68ugRXbGQVF3tU5KIXSP6+glvjfWezufNAFemvf40koZGMJSntU2GAjEKBzt', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', 'PyVfXw0bdcapZFaThsv7kDRo1osPEmI5gNvH1QzutsHsPeKTdjuEk6WAWg==', '', 'student', 'STU001', TRUE),
('yrAADhR5LXgO1GigARLKH7wexjVJ9pu8x2enMl7xZytKxFfvrxHEo//99F1FTw4LlMhS1Hc=', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', '11mJogmd6fJjLjVNzM4bOIRHGLeEwAZrRN9s6EnvkRQNSuqe', 'lPLe2VaLPk8you1Y654KM0S7a4zKKTGjvJjjIHSf7GWCx2VyCIk=', 'student', 'STU002', TRUE),
('4litX1/vVP+rR7zySgh/kTgUWC9tr1EY0ltkoGj4SFD/0XdbFLjZtHbupTuIu+kfou9pZAWE', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', 'Xkz86ZGnvbOWUNDF+5WeuR7wrkOKHY0w6Eigqh14bTAYSQzR3sA=', '2Yvd4/+z6EyjCLRP0sExK40/zrMRALnQiSF0qqCyANK5G1DtiU8922A=', 'student', 'STU003', TRUE),
('CLUXIzxC+SI1tdvFzN5uB1YTGTDHUaUnZsJo6zXnVM3XkHdeWXp0LAu9TE30GvA+pg==', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', 'Jjngnw4pO8Qo3dInJv6XOrOrVlFzU2uNInnwD2xUyINIeHw=', 'vMZF2PDEdZw4JdLwZD/b/V68gfVz2husnCUpSSGLm1c=', 'student', 'STU004', TRUE),
('3kJvWthICg0eZXp/MpQsiIxRJb1m2y+VTHvOHtI9nObVC6myE5DoVPfkbk/LS9h4TJ2COjw=', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', 'OoJnvLpSZFUb3tNB3DBs0EjpQRyO4dBPJLGr03e3QVkQZgV9', 'OCIlMg05eaLLkJ7UvwgG10lSepEVzwSwor0oQiGCtVcugXwGNQ==', 'student', 'STU005', TRUE),
('Wvhd+vWFB/JgT5U6loKtaFBPXUFJg91LpSRoJtxcVnHAEbp/D3LvY3O4qpGK9wV6r1Y6', '$2a$10$1gWmkyIBZegPgg7gAJzSveSm3iXQh2o2H7EXAV8jzGi7lUo4yPd7C', 'kdwgL+lC6/t94RiSrq6j7ZypZq2yXklbofHqqWghwD9G23aDsEU=', 'UmJFT/evK0+SWCkoQgaHrMGLCRnHUtY2btGbCmvqJZ7k9tui10Gltw==', 'student', 'STU006', TRUE);
