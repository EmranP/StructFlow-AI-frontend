import { HttpClient, HttpParams } from '@angular/common/http'
import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import { environment } from '../../environments/environment'
import { CreateProjectDto, ProjectItem, Projects } from '../models'

@Injectable({ providedIn: 'root' })
export class ProjectService {
	private apiUrl = `${environment.apiUrl}/api`

	constructor(private http: HttpClient) {}

	getAll(page = 1, limit = 10): Observable<Projects> {
		const params = new HttpParams()
			.set('page', page.toString())
			.set('limit', limit.toString())

		return this.http.get<Projects>(`${this.apiUrl}/project/all`, { params })
	}

	getById(id: string): Observable<ProjectItem> {
		return this.http.get<ProjectItem>(`${this.apiUrl}/project/${id}`)
	}

	create(dto: CreateProjectDto): Observable<{ message: string; id: string }> {
		return this.http.post<{ message: string; id: string }>(
			`${this.apiUrl}/project/new`,
			dto
		)
	}

	edit(
		id: string,
		dto: Partial<CreateProjectDto>
	): Observable<{ message: string }> {
		return this.http.patch<{ message: string }>(
			`${this.apiUrl}/project/edit/${id}`,
			dto
		)
	}

	delete(id: string): Observable<{ message: string }> {
		return this.http.delete<{ message: string }>(
			`${this.apiUrl}/project/remove/${id}`
		)
	}

	startGeneration(
		projectId: string,
		model: string
	): Observable<{ id: string; status: string }> {
		return this.http.post<{ id: string; status: string }>(
			`${this.apiUrl}/project/gen/${projectId}`,
			{
				model,
			}
		)
	}
}
