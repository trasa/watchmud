-- psql -U watchmud

CREATE TABLE players (
  player_name varchar(100) primary key not null,
  current_health integer not null,
  max_health integer not null,
  race integer not null,
  class integer not null
);

CREATE TABLE player_inventory (
  player_name varchar(1000) not null REFERENCES players(player_name),
  instance_id UUID not null,
  zone_id varchar(1000) not null,
  definition_id varchar(1000) not null,
  primary key (player_name, instance_id)
);



GRANT SELECT, INSERT, UPDATE ON TABLE players TO watchmud;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE player_inventory to watchmud;

INSERT INTO players (player_name, current_health, max_health, race, class)
values ('somedood', 100, 100, 0, 0);
