package rules

// BaseArmorClass is what anything is worth wearing nothing at all, and the
// bottom of the only scale combat has.
//
// It lives here rather than next to the equipment that adds to it because
// both sides of an attack are measured against it: a player's armor class is
// this plus what they are wearing (object.Equipment.ArmorClass), and a mob's
// is this unless its definition says otherwise. combat.AttemptMeleeAttack
// compares a d20 to whichever of the two is defending, so the two have to
// mean the same thing. They did not when this constant was private to
// object/: mobs were whatever number content typed, and a mob whose file
// forgot to type one was AC 0 -- which, against a die that cannot roll below
// 1, is a monster that can never be missed.
const BaseArmorClass = 10
