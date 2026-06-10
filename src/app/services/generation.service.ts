import { HttpClient, HttpParams } from '@angular/common/http'
import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import { environment } from '../../environments/environment'
import { Generation, GenerationResponse, Template } from '../models'

@Injectable({ providedIn: 'root' })
export class GenerationService {
	private apiUrl = `${environment.apiUrl}/api`

	constructor(private http: HttpClient) {}

	getAll(
		projectId: string,
		page = 1,
		limit = 10
	): Observable<GenerationResponse> {
		const params = new HttpParams()
			.set('page', page.toString())
			.set('limit', limit.toString())
		return this.http.get<GenerationResponse>(
			`${this.apiUrl}/gen/all/${projectId}`,
			{
				params,
			}
		)
	}

	getById(id: string): Observable<Generation> {
		return this.http.get<Generation>(`${this.apiUrl}/gen/${id}`)
	}

	getTemplates(generationId: string): Observable<Template[]> {
		return this.http.get<Template[]>(
			`${this.apiUrl}/gen/${generationId}/templates`
		)
	}

	getDownloadUrl(generationId: string): string {
		return `${this.apiUrl}/gen/download/${generationId}`
	}
}
