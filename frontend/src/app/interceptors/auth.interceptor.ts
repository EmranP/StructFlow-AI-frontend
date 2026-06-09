import { HttpClient, HttpInterceptorFn } from '@angular/common/http'
import { inject } from '@angular/core'
import { catchError, switchMap, throwError } from 'rxjs'
import { AuthResponse } from '../models'
import { AuthService } from '../services/auth.service'

export const authInterceptor: HttpInterceptorFn = (req, next) => {
	const authService = inject(AuthService)
	const token = authService.getToken()

	if (token) {
		const clonedReq = req.clone({
			setHeaders: {
				Authorization: `Bearer ${token}`,
			},
		})
		return next(clonedReq)
	}

	return next(req)
}

export const refreshInterceptor: HttpInterceptorFn = (req, next) => {
	if (req.url.includes('/auth/')) return next(req)

	return next(req).pipe(
		catchError(err => {
			if (err.status !== 401) return throwError(() => err)

			return inject(HttpClient)
				.get<AuthResponse>('http://localhost:3000/api/auth/refresh')
				.pipe(
					switchMap(res => {
						localStorage.setItem('token', res.accessToken)
						const retried = req.clone({
							setHeaders: { Authorization: `Bearer ${res.accessToken}` },
						})
						return next(retried)
					}),
					catchError(() => {
						localStorage.removeItem('token')
						window.location.href = '/auth'
						return throwError(() => err)
					})
				)
		})
	)
}
