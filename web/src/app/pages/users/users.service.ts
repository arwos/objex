import { Injectable } from '@angular/core';
import { RequestService } from '@uxwb/services';

@Injectable({
  providedIn: 'root',
})
export class UsersService {

  constructor(private readonly http: RequestService) { }

}
