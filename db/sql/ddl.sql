CREATE TABLE players (
  player_name varchar(100) primary key not null,
  current_health integer not null,
  max_health integer not null
);

INSERT INTO players (player_name, current_health, max_health) 
values ('somedood', 100, 100);
