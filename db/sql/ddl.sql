-- psql -U watchmud

CREATE TABLE players (
  player_id serial primary key not null,
  player_name varchar(100) not null unique,
  current_health integer not null,
  max_health integer not null,
  race integer not null,
  class integer not null,
  last_zone_id varchar(100),
  last_room_id varchar(100),
  slots JSONB null
);

CREATE TABLE player_inventory (
  player_id integer not null REFERENCES players(player_id),
  instance_id UUID not null,
  zone_id varchar(1000) not null,
  definition_id varchar(1000) not null,
  primary key (player_id, instance_id)
);


GRANT SELECT, INSERT, UPDATE ON TABLE players TO watchmud;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE player_inventory to watchmud;

INSERT INTO players (player_name, current_health, max_health, race, class)
values ('somedood', 100, 100, 0, 0);
