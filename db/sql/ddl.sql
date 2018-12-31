CREATE TABLE players (
  player_name varchar(100) primary key not null,
  current_health integer not null,
  max_health integer not null,
  race integer not null,
  class integer not null
);

GRANT SELECT, INSERT, UPDATE ON TABLE players TO watchmud;

INSERT INTO players (player_name, current_health, max_health, race, class)
values ('somedood', 100, 100, 0, 0);
