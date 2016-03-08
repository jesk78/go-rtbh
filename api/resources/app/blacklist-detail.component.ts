import {Component, OnInit}  from 'angular2/core';
import {BlacklistEntry}     from './blacklist-entry';
import {BlacklistService}   from './blacklist.service';

@Component({
    selector: 'blacklist-details',
    inputs: ['entry'],
    template: `
    <div *ngIf="selectedEntry">
    <h3>Details for {{selectedEntry.address}}</h3>
    <i>Reason</i>: {{selectedEntry.reason}}
    </div>
    `
})

export class BlacklistDetailComponent implements OnInit {
    selectedEntry: BlacklistEntry;

    constructor(private _blacklistService: BlacklistService) {}

    ngOnInit() {
        this._blacklistService.selectedEntry$.subscribe(entry => this.selectedEntry = entry);
    }
}
