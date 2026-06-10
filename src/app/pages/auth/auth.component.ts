import { CommonModule } from '@angular/common'
import { Component, inject } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { Router } from '@angular/router'
import { AuthService } from '../../services/auth.service'

@Component({
	selector: 'app-auth',
	standalone: true,
	imports: [CommonModule, FormsModule],
	templateUrl: './auth.component.html',
	styleUrls: ['./auth.component.scss'],
})
export class AuthComponent {
	auth = inject(AuthService)
	router = inject(Router)

	mode: 'login' | 'register' = 'login'
	email = ''
	password = ''
	loading = false
	error = ''

	resent = false

	toggle() {
		this.mode = this.mode === 'login' ? 'register' : 'login'
		this.error = ''
		this.resent = false
	}

	submit() {
		if (!this.email || !this.password) {
			this.error = 'Fill all fields'
			return
		}
		this.loading = true
		this.error = ''
		this.resent = false

		if (this.mode === 'login') {
			this.auth.login(this.email, this.password).subscribe({
				next: () => this.router.navigate(['/projects']),
				error: e => {
					this.loading = false
					this.error = e?.error?.message || 'Invalid credentials'
				},
			})
		} else {
			this.auth.register(this.email, this.password).subscribe({
				next: () => {
					// Успешная регистрация — идём на verify
					this.loading = false
					this.router.navigate(['/verify-email'], {
						queryParams: { email: this.email },
					})
				},
				error: e => {
					this.loading = false
					const msg = e?.error?.message || ''

					if (e?.status === 409 && msg.includes('resent')) {
						this.router.navigate(['/verify-email'], {
							queryParams: { email: this.email },
						})
						return
					}

					this.error = msg || 'Registration failed'
				},
			})
		}
	}
}
