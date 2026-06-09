import { HttpClient } from '@angular/common/http'
import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import { ModelResponse } from '../models'

@Injectable({ providedIn: 'root' })
export class AiService {
	private apiUrl = 'http://localhost:3000/api'

	constructor(private http: HttpClient) {}

	getModels(): Observable<ModelResponse[]> {
		return this.http.get<ModelResponse[]>(`${this.apiUrl}/ai/models`)
	}
}
