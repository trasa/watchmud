import { Component, OnInit, Input } from '@angular/core';
import { Zone } from '../zone';

@Component({
  selector: 'app-zone-detail',
  templateUrl: './zone-detail.component.html',
  styleUrls: ['./zone-detail.component.scss']
})
export class ZoneDetailComponent implements OnInit {
  @Input() zone: Zone;

  constructor() { }

  ngOnInit() {
  }

}
