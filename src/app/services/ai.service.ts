import { HttpClient } from '@angular/common/http'
import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import { environment } from '../../environments/environment'
import { ModelResponse } from '../models'

@Injectable({ providedIn: 'root' })
export class AiService {
	constructor(private http: HttpClient) {}

	getModels(): Observable<ModelResponse[]> {
		return this.http.get<ModelResponse[]>(`${environment.apiUrl}/api/ai/models`)
	}
}
