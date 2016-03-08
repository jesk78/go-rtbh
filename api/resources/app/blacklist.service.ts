import {Injectable} from 'angular2/core';
import {Http, Response} from 'angular2/http';
import {BlacklistEntry} from './blacklist-entry';
import {Observable}     from 'rxjs/Observable';

@Injectable()
export class BlacklistService {
    selectedEntry$: Observable<BlacklistEntry>;
    private _selectedEntryObserver: any;

    constructor (private http: Http) {
        this.selectedEntry$ = new Observable(
            observer => this._selectedEntryObserver = observer
        ).share();
    };

    private _blacklist_path = '/v1/blacklist';

    getEntries() {
        return this.http.get(this._blacklist_path)
            .map(res => <BlacklistEntry[]> res.json())
            .catch(this.apiError);
    }

    setSelectedEntry(entry: BlacklistEntry) {
        this._selectedEntryObserver.next(entry);
    }

    private apiError(error: Response) {
        console.log(error);

        return Observable.throw(error.json().error || 'Server Error');
    }
}
