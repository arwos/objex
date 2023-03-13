import { Injectable } from '@angular/core';
import { RequestService } from '@uxwb/services';
import { map, Observable, take } from 'rxjs';
import { Group } from 'src/app/pages/groups/models';

@Injectable({
  providedIn: 'root',
})
export class GroupsService {

  constructor(private readonly http: RequestService) { }

  list(): Observable<Group[]> {
    return this.http.get('groups/list')
      .pipe(take(1), map((data) => <Group[]>data));
  }
}
