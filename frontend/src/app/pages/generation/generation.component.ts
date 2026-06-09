import { CommonModule } from '@angular/common'
import { Component, inject, OnDestroy, OnInit } from '@angular/core'
import { ActivatedRoute, Router, RouterLink } from '@angular/router'
import { LoadingWaveComponent } from '../../components/loading-wave/loading-wave.component'
import { Generation } from '../../models'
import { AuthService } from '../../services/auth.service'
import { GenerationService } from '../../services/generation.service'

@Component({
	selector: 'app-generation',
	standalone: true,
	imports: [CommonModule, RouterLink, LoadingWaveComponent],
	templateUrl: './generation.component.html',
	styleUrls: ['./generation.component.scss'],
})
export class GenerationComponent implements OnInit, OnDestroy {
	route = inject(ActivatedRoute)
	router = inject(Router)
	generationService = inject(GenerationService)
	authService = inject(AuthService)

	generation: Generation | null = null
	genId = ''
	pollInterval: any
	countdown = 5
	countdownInterval: any
	error = ''
	showBar = false

	ngOnInit() {
		this.authService.getMe().subscribe({
			error: err => {
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
		this.genId = this.route.snapshot.paramMap.get('id')!
		this.checkStatus()
		this.startPolling()
	}

	ngOnDestroy() {
		this.stopPolling()
	}

	checkStatus() {
		this.generationService.getById(this.genId).subscribe({
			next: gen => {
				this.generation = gen
				setTimeout(() => {
					this.showBar = true
				}, 0)
				if (gen.status === 'completed' || gen.status === 'failed') {
					this.stopPolling()
				}
			},
			error: err => {
				this.error = 'Failed to check status'
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
	}

	startPolling() {
		this.countdown = 20

		this.countdownInterval = setInterval(() => {
			this.countdown--
			if (this.countdown <= 0) this.countdown = 20
		}, 1000)

		this.pollInterval = setInterval(() => {
			if (
				!this.generation ||
				(this.generation.status !== 'completed' &&
					this.generation.status !== 'failed')
			) {
				this.checkStatus()
				this.countdown = 20
			}
		}, 30000)
	}

	stopPolling() {
		if (this.pollInterval) clearInterval(this.pollInterval)
		if (this.countdownInterval) clearInterval(this.countdownInterval)
	}

	viewStructure() {
		this.router.navigate(['/structure', this.genId])
	}

	getStatusLabel(s: string) {
		if (s === 'pending') return 'Waiting in queue...'
		if (s === 'processing') return 'AI is designing your structure...'
		if (s === 'completed') return 'Generation complete!'
		if (s === 'failed') return 'Generation failed'
		return s
	}
}
