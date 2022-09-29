CREATE TABLE quotes (
    uuidkey varchar(36) PRIMARY KEY,
    quote varchar(1500) NOT NULL,
    author varchar(50) NOT NULL,
);

INSERT INTO quotes (uuidkey, quote, author)
VALUES ('4589f6c5-4129-4919-b3db-de45edee244b', 'Reflection is never clear.', 'Joe Burrow'),
('adfec432-e3ca-4a18-9ce3-0db41ae67784', 'Don''t just check errors handle them gracefully.' 'Oprah'),
('3736cfb7-1569-4cea-b9d7-cc438afa2ab8', 'A little copying is better than a little dependency.', 'Vienna Erhart'),
('2c9d9b91-e829-4495-8e7c-4991ddc8c5d6', 'The bigger the interface, the weaker the abstraction.', 'Josh Rose'),
('cca8abd8-37f8-4ed9-8853-2cfd734d6965',  'Don''t panic.', 'Queen of England');
