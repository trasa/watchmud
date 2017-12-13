import { Component, OnInit } from '@angular/core';
import { Zone } from '../zone';

@Component({
  selector: 'app-zones',
  templateUrl: './zones.component.html',
  styleUrls: ['./zones.component.css']
})
export class ZonesComponent implements OnInit {
  zone: Zone = {
    id: 'sample',
    name: 'Sample Zone',
    resetMode: 2,
    lifetimeMinutes: 3,
    enabled: true
  };


  constructor() { }

  ngOnInit() {
  }

}
