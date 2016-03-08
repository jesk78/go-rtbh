import {Component}          from 'angular2/core';
import {HTTP_PROVIDERS}     from 'angular2/http';
import {RouteConfig, ROUTER_DIRECTIVES, ROUTER_PROVIDERS} from 'angular2/router';

import {BlacklistService}   from './blacklist.service';
import {BlacklistComponent} from './blacklist.component';
import {DashboardComponent} from './dashboard.component';

@Component({
    selector: 'rtbh-app',
    template: `
    <nav class="navbar navbar-default">
      <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="#">{{title}}</a>
        </div>
		<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">

		  <ul class="nav navbar-nav">
            <li><a href="#" [routerLink]="['Dashboard']">Dashboard</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="#" [routerLink]="['Blacklist']">Blacklist</a></li>
          </ul>
        </div><!-- /.navbar-collapse -->
      </div><!-- /.container-fluid -->
    </nav>
    <router-outlet></router-outlet>
    `,
    directives: [ROUTER_DIRECTIVES],
    providers: [
        ROUTER_PROVIDERS,
        HTTP_PROVIDERS,
        BlacklistService,
    ]
})

@RouteConfig([
    {
        path: '/blacklist',
        name: 'Blacklist',
        component: BlacklistComponent
    },
    {
        path: '/dashboard',
        name: 'Dashboard',
        component: DashboardComponent
    }
])

export class AppComponent {
    title = 's/RTBH Blacklist Console';
}
