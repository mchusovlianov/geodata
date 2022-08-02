INSERT INTO countries (uuid, `code`, `name`, `date_created`, `date_updated`) VALUES
     ('77eabf6e-30a8-44d0-8952-029d2ca06872', 'AL', 'Alabnia', '2021-01-01 00:00:01.000001+00', '2021-01-01 00:00:01.000001+00');

INSERT INTO locations (`uuid`, `city_uuid`, `ip`, `mystery_value`, `latitude`, `longitude`, `date_created`, `date_updated`) VALUES
      ('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '32.123.12.2', 1232412415, -50.023, 54.321, '2021-01-01 00:00:01.000001+00', '2021-01-01 00:00:01.000001+00');

INSERT INTO cities (uuid, `country_uuid`, `name`, date_created, date_updated) VALUES
      ('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '77eabf6e-30a8-44d0-8952-029d2ca06872', 'Test city #1', '2021-01-01 00:00:01.000001+00', '2021-01-01 00:00:01.000001+00');
