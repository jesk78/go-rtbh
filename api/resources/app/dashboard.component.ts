import {Component}      from 'angular2/core';

@Component({
    selector: 'dashboard',
    template: `
    <h3>{{title}}</h3>
    `
})

export class DashboardComponent {
    title = 'Dashboard';
}
