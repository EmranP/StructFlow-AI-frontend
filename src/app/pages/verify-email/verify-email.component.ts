import { CommonModule } from '@angular/common'
import { Component, inject, OnDestroy, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { ActivatedRoute, Router } from '@angular/router'
import { AuthService } from '../../services/auth.service'

@Component({
	selector: 'app-verify-email',
	standalone: true,
	imports: [CommonModule, FormsModule],
	templateUrl: './verify-email.component.html',
	styleUrls: ['./verify-email.component.scss'],
})
export class VerifyEmailComponent implements OnInit, OnDestroy {
	auth = inject(AuthService)
	route = inject(ActivatedRoute)
	router = inject(Router)

	email = ''
	code = ''
	loading = false
	resending = false
	error = ''
	successMsg = ''

	// Cooldown таймер для повторной отправки
	cooldown = 0
	private cooldownInterval: any

	ngOnInit() {
		this.email = this.route.snapshot.queryParamMap.get('email') || ''
		this.startCooldown(60) // первая отправка была только что
	}

	ngOnDestroy() {
		if (this.cooldownInterval) clearInterval(this.cooldownInterval)
	}

	startCooldown(seconds: number) {
		this.cooldown = seconds
		this.cooldownInterval = setInterval(() => {
			this.cooldown--
			if (this.cooldown <= 0) {
				clearInterval(this.cooldownInterval)
				this.cooldown = 0
			}
		}, 1000)
	}

	verify() {
		if (!this.code.trim()) {
			this.error = 'Enter the verification code'
			return
		}
		this.loading = true
		this.error = ''

		this.auth.verifyEmail({ email: this.email, code: this.code }).subscribe({
			next: () => {
				this.loading = false
				this.router.navigate(['/projects'])
			},
			error: e => {
				this.loading = false
				this.error = e?.error?.message || 'Invalid or expired code'
			},
		})
	}

	resend() {
		if (this.cooldown > 0 || this.resending) return
		this.resending = true
		this.error = ''
		this.successMsg = ''

		this.auth.register(this.email, '').subscribe({
			next: () => {
				this.resending = false
				this.successMsg = 'Code sent!'
				this.startCooldown(60)
			},
			error: e => {
				this.resending = false

				if (e?.status === 202 || e?.error?.message?.includes('resent')) {
					this.successMsg = 'Code resent to your email'
					this.startCooldown(60)
					return
				}
				this.error = e?.error?.message || 'Failed to resend'
			},
		})
	}
}
