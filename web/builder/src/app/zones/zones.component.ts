import { Component, OnInit } from '@angular/core';
import { Zone } from '../zone';
import { ZONES } from '../mock-zones';

@Component({
  selector: 'app-zones',
  templateUrl: './zones.component.html',
  styleUrls: ['./zones.component.css']
})
export class ZonesComponent implements OnInit {
  zones = ZONES;

  selectedZone: Zone;

  onSelect(zone: Zone): void {
    this.selectedZone = zone;
  }


  constructor() { }

  ngOnInit() {
  }

}
