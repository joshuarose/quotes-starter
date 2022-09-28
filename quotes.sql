CREATE TABLE quotes (
Id varchar PRIMARY KEY,
  Phrase varchar,
  Author varchar NOT NULL
);

INSERT 
    INTO quotes (Id, Phrase, Author)
VALUES ('374be3f1-956a-4169-874a-0632c09a2599', 'Don't communicate by sharing memory, share memory by communicating','Rob Pike'),
('a4539044-da8d-4064-bb05-2421abd4c77d', 'With the unsafe package there are no guarantees.', 'Rob Pike'),
('068faa87-9afa-4f7f-8aed-ff2d303c79e5', 'A little copying is better than a little dependency.', 'Rob Pike'),
('0f4036b0-d49a-46b9-9ec2-577fbfd4f714', 'Design the architecture, name the components, document the details.', 'Rob Pike'),
('10a2781c-113f-4c49-a670-8ed322882f1a', 'Don't just check errors, handle them gracefully.', 'Rob Pike'),
('77efbc8b-2289-45ee-9461-b1f602fecf3e', 'Avoid unused method receiver names', 'Kalese Carpenter'),
('211cf4f3-3893-43b8-a1d2-88aedc14df5a', 'Gofmt's style is no one's favorite, yet gofmt is everyone's favorite', 'Rob Pike'),
('323d8e20-7975-4ff1-af6d-99dc7f57f35a', 'For brands or words with more than 1 capital letter, lowercase all letters', 'Kalese Carpenter');
