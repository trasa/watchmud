import { Component, OnInit } from '@angular/core';
import {Zone} from '../zone';
import { ZoneService } from '../zone.service';

@Component({
  selector: 'ngx-zones',
  templateUrl: './zones.component.html',
  styleUrls: ['./zones.component.scss'],
})
export class ZonesComponent implements OnInit {
  zones: Zone[];
  selectedZone: Zone;

  zone: Zone = {
    id: 'sample',
    name: 'Sample Zone',
    resetMode: 0,
    lifetimeMinutes: 0,
    enabled: true,
  };
  constructor(private zoneService: ZoneService) { }

  ngOnInit() {
    this.loadZones();
  }

  loadZones(): void {
    this.zones = this.zoneService.getZones();
  }

  onSelect(zone: Zone): void {
    this.selectedZone = zone;
  }

}
