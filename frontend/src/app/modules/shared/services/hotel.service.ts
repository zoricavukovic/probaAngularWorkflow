import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {ConfigService} from "./config-service/config.service";

@Injectable({
  providedIn: 'root',
})
export class HotelService {
  constructor(
    private http: HttpClient,
    private configService: ConfigService
  ) {}

  hello(): Observable<string> {
    return this.http.get(`/auth/health`, { responseType: 'text' }
    );
  }
}
