import { HttpClient } from '@angular/common/http'
import { Injectable, signal } from '@angular/core'
import { Observable, tap } from 'rxjs'
import { environment } from '../../environments/environment'
import { AuthResponse, GetID, MessageResponse, VerifyEmailDto } from '../models'

@Injectable({ providedIn: 'root' })
export class AuthService {
	private apiUrl = `${environment.apiUrl}/api`
	isLoggedIn = signal<boolean>(this.hasToken())

	constructor(private http: HttpClient) {}

	private hasToken(): boolean {
		return !!localStorage.getItem('token')
	}

	getToken(): string | null {
		return localStorage.getItem('token')
	}

	getMe(): Observable<GetID> {
		return this.http.get<GetID>(`${this.apiUrl}/auth/me`)
	}

	register(email: string, password: string): Observable<MessageResponse> {
		return this.http.post<MessageResponse>(`${this.apiUrl}/auth/register`, {
			email,
			password,
		})
	}

	login(email: string, password: string): Observable<AuthResponse> {
		return this.http
			.post<AuthResponse>(`${this.apiUrl}/auth/login`, { email, password })
			.pipe(
				tap(res => {
					localStorage.setItem('token', res.accessToken)
					this.isLoggedIn.set(true)
				})
			)
	}

	verifyEmail(dto: VerifyEmailDto): Observable<AuthResponse> {
		return this.http
			.post<AuthResponse>(`${this.apiUrl}/auth/verify-email`, dto)
			.pipe(
				tap(res => {
					localStorage.setItem('token', res.accessToken)
					this.isLoggedIn.set(true)
				})
			)
	}

	refresh(): Observable<AuthResponse> {
		return this.http.get<AuthResponse>(`${this.apiUrl}/auth/refresh`).pipe(
			tap(res => {
				localStorage.setItem('token', res.accessToken)
			})
		)
	}

	logout(): Observable<MessageResponse> {
		return this.http
			.post<MessageResponse>(`${this.apiUrl}/auth/logout`, {})
			.pipe(
				tap(() => {
					localStorage.removeItem('token')
					this.isLoggedIn.set(false)
				})
			)
	}
}
