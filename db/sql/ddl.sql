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
  str_bonus integer not null default 0,
  dex_bonus integer not null default 0,
  con_bonus integer not null default 0,
  int_bonus integer not null default 0,
  wis_bonus integer not null default 0,
  cha_bonus integer not null default 0
);

INSERT INTO races (race_id, race_group_id, race_name, str_bonus, dex_bonus, con_bonus, int_bonus, wis_bonus, cha_bonus) values
                  (0, 0, 'Human', 1, 1, 1, 1, 1, 1);
INSERT INTO races (race_id, race_group_id, race_name, con_bonus, wis_bonus) values
                  (1, 1, 'Hill Dwarf', 2, 1);
INSERT INTO races (race_id, race_group_id, race_name, con_bonus, str_bonus) values
                  (2, 1, 'Mountain Dwarf', 2, 2);
INSERT INTO races (race_id, race_group_id, race_name, dex_bonus, int_bonus) values
                  (3, 2, 'High Elf', 2, 1);
INSERT INTO races (race_id, race_group_id, race_name, dex_bonus, wis_bonus) values
                  (4, 2, 'Wood Elf', 2, 1);
INSERT INTO races (race_id, race_group_id, race_name, dex_bonus, cha_bonus) values
                  (5, 3, 'Lightfoot Halfling', 2, 1);
INSERT INTO races (race_id, race_group_id, race_name, dex_bonus, con_bonus) values
                  (6, 3, 'Stout Halfling', 2, 1);

CREATE TABLE classes (
    class_id integer primary key not null,
    class_name varchar(100) not null,
    ability_preference JSONB not null
);

INSERT INTO classes (class_id, class_name, ability_preference) values (0, 'Fighter', '{ "a": ["str", "dex", "con"]}');
INSERT INTO classes (class_id, class_name, ability_preference) values (1, 'Cleric', '{ "a": ["wis", "con", "cha"]}');
INSERT INTO classes (class_id, class_name, ability_preference) values (2, 'Rogue', '{ "a": ["dex", "str", "con"]}');
INSERT INTO classes (class_id, class_name, ability_preference) values (3, 'Wizard', '{ "a": ["int", "wis", "con"]}');


CREATE TABLE players (
  player_id serial primary key,
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

CREATE TABLE player_inventory
(
    inventory_id  serial primary key,
    player_id     integer       not null REFERENCES players (player_id),
    instance_id   UUID          not null,
    zone_id       varchar(1000) not null,
    definition_id varchar(1000) not null,
    unique (player_id, instance_id)
);



-- INSERT INTO players (player_name, current_health, max_health, race_id, class, strength, dexterity, constitution, intelligence, wisdom, charisma)
-- values ('somedood', 100, 100, 0, 0, 15, 10, 14, 8, 13, 12);


GRANT SELECT, INSERT, UPDATE ON TABLE players TO watchmud;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE player_inventory to watchmud;
GRANT SELECT ON TABLE race_group to watchmud;
GRANT SELECT ON TABLE races to watchmud;
GRANT SELECT ON TABLE classes to watchmud;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public to watchmud;