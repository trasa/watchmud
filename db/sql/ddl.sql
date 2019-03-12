-- psql -U watchmud

CREATE TABLE race_group (
  race_group_id integer primary key not null,
  race_group_name varchar(100) not null
);

INSERT INTO race_group (race_group_id, race_group_name) values (0, 'Human');
INSERT INTO race_group (race_group_id, race_group_name) values (1, 'Dwarf');
INSERT INTO race_group (race_group_id, race_group_name) values (2, 'Elf');
INSERT INTO race_group (race_group_id, race_group_name) values (3, 'Halfling');

CREATE TABLE races (
  race_id integer primary key not null,
  race_group_id integer not null references race_group(race_group_id),
  race_name varchar(100) not null,
  ability_scores JSONB null
);

INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (0, 0, 'Human', '[{"str": 1}, {"dex": 1}, {"con": 1}, {"int": 1}, {"wis": 1}, {"cha": 1} ]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (1, 1, 'Hill Dwarf', '[ {"con": 2}, {"wis": 1} ]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (2, 1, 'Mountain Dwarf', '[ {"con": 2}, {"str": 2}]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (3, 2, 'High Elf', '[ {"dex": 2}, {"int": 1}]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (4, 2, 'Wood Elf', '[ {"dex": 2}, {"wis": 1}]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (5, 3, 'Lightfoot Halfling', '[ {"dex": 2}, {"cha": 1}]');
INSERT INTO races (race_id, race_group_id, race_name, ability_scores) values (6, 3, 'Stout Halfling', '[{"dex": 2}, {"con": 1}]');


CREATE TABLE players (
  player_id serial primary key not null,
  player_name varchar(100) not null unique,
  current_health integer not null,
  max_health integer not null,
  race_id integer not null references races(race_id),
  class integer not null,
  last_zone_id varchar(100),
  last_room_id varchar(100),
  strength integer not null,
  dexterity integer not null,
  constitution integer not null,
  intelligence integer not null,
  wisdom integer not null,
  charisma integer not null,
  slots JSONB null
);

CREATE TABLE player_inventory (
  player_id integer not null REFERENCES players(player_id),
  instance_id UUID not null,
  zone_id varchar(1000) not null,
  definition_id varchar(1000) not null,
  primary key (player_id, instance_id)
);



INSERT INTO players (player_name, current_health, max_health, race_id, class, strength, dexterity, constitution, intelligence, wisdom, charisma)
values ('somedood', 100, 100, 0, 0, 15, 10, 14, 8, 13, 12);


GRANT SELECT, INSERT, UPDATE ON TABLE players TO watchmud;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE player_inventory to watchmud;
GRANT SELECT ON TABLE race_group to watchmud;
GRANT SELECT ON TABLE races to watchmud;
