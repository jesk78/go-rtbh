import {Component}      from 'angular2/core';
import {BlacklistEntry} from './blacklist-entry';

@Component({
    selector: 'blacklist-details',
    inputs: ['entry'],
    template: `
    <div *ngIf="entry">
    <h3>Details for {{entry.address}}</h3>
    <i>Reason</i>: {{entry.reason}}
    </div>
    `
})

export class BlacklistDetailComponent {
}
