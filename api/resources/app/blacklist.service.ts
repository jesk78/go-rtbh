import {ENTRIES} from './mock-blacklist';
import {Injectable} from 'angular2/core';

@Injectable()
export class BlacklistService {
    getEntries() {
        return Promise.resolve(ENTRIES);
    }
}
