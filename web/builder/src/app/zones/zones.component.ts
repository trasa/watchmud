import { Component, OnInit } from '@angular/core';
import {Zone} from "../zone";
import {ZONES} from "../mock-zones";

@Component({
  selector: 'app-zones',
  templateUrl: './zones.component.html',
  styleUrls: ['./zones.component.scss']
})
export class ZonesComponent implements OnInit {
  zones = ZONES;
  selectedZone: Zone;

  zone: Zone = {
    id: "sample",
    name: "Sample Zone",
    resetMode: 0,
    lifetimeMinutes: 0,
    enabled: true,
  };
  constructor() { }

  ngOnInit() {
  }

  onSelect(zone: Zone): void {
    this.selectedZone = zone;
  }

}
