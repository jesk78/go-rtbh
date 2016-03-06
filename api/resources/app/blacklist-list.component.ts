import {Component}          from 'angular2/core';
import {BlacklistEntry}     from './blacklist-entry';
import {BlacklistService}   from './blacklist.service';

@Component({
    selector: 'blacklist-list',
    inputs: ['entries'],
    template: `
    <h3>Blacklist entries</h3>
    <ul class="list-group">
        <li class="list-group-item" *ngFor="#entry of entries"
            [class.selected]="entry === selectedEntry"
            (click)="onSelect(entry)">
            {{entry.address}}
        </li>
    </ul>
    `
})

export class BlacklistListComponent {
    selectedEntry: BlacklistEntry;
    entries: BlacklistEntry[];

    onSelect(entry) {
        this.selectedEntry = entry;
    }
}
