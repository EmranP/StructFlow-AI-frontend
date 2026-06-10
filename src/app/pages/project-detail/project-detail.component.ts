import { CommonModule } from '@angular/common'
import { Component, inject, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { ActivatedRoute, Router, RouterLink } from '@angular/router'
import { Generation, ModelResponse, ProjectItem } from '../../models'
import { AiService } from '../../services/ai.service'
import { AuthService } from '../../services/auth.service'
import { GenerationService } from '../../services/generation.service'
import { ProjectService } from '../../services/project.service'

@Component({
	selector: 'app-project-detail',
	standalone: true,
	imports: [CommonModule, RouterLink, FormsModule],
	templateUrl: './project-detail.component.html',
	styleUrls: [
		'./project-detail.component.scss',
		'./project-detail-icon.component.scss',
		'../../app.component.scss',
	],
})
export class ProjectDetailComponent implements OnInit {
	route = inject(ActivatedRoute)
	router = inject(Router)
	projectService = inject(ProjectService)
	generationService = inject(GenerationService)
	authService = inject(AuthService)
	aiService = inject(AiService)

	project: ProjectItem | null = null
	generations: Generation[] = []
	loading = true
	gensLoading = false
	generating = false
	editing = false
	saving = false
	modelsLoading = true

	models: ModelResponse[] = []
	selectedModel = 'gemini'

	// Pagination
	page = 1
	limit = 10
	totalCount = 0

	editForm: any = {}

	readonly modelConfig: Record<
		string,
		{ icon: string; color: string; bg: string }
	> = {
		gemini: {
			icon: '../../../assets/gemini.svg',
			color: '#4285f4',
			bg: 'rgba(66, 133, 244, 0.12)',
		},
		claude: {
			icon: '../../../assets/claude-ai-icon.svg',
			color: '#d97706',
			bg: 'rgba(217, 119, 6, 0.12)',
		},
		gpt: {
			icon: '../../../assets/openai.svg',
			color: '#10a37f',
			bg: 'rgba(16, 163, 127, 0.12)',
		},
	}

	get allModelsUnavailable(): boolean {
		return (
			!this.modelsLoading &&
			(this.models.length === 0 || this.models.every(m => !m.available))
		)
	}

	get totalPages(): number {
		return Math.ceil(this.totalCount / this.limit)
	}

	get pages(): number[] {
		return Array.from({ length: this.totalPages }, (_, i) => i + 1)
	}

	ngOnInit() {
		this.authService.getMe().subscribe({
			error: err => {
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
		const id = this.route.snapshot.paramMap.get('id')!
		this.loadData(id)
		this.loadModels()
	}

	loadModels() {
		this.modelsLoading = true
		this.aiService.getModels().subscribe({
			next: models => {
				this.models = models
				if (
					this.models.length &&
					!this.models.find(m => m.id === this.selectedModel)
				) {
					this.selectedModel = this.models[0].id
				}
				this.modelsLoading = false
			},
			error: () => {
				// Fallback если /ai/models недоступен
				this.models = [
					{
						id: 'gemini',
						name: 'Gemini',
						available: false,
					},
					{
						id: 'claude',
						name: 'Claude',
						available: false,
					},
					{
						id: 'gpt',
						name: 'GPT-4',
						available: false,
					},
				]
				this.modelsLoading = false
			},
		})
	}

	loadData(id: string) {
		this.loading = true
		this.projectService.getById(id).subscribe({
			next: p => {
				this.project = p
				this.editForm = { ...p }
				this.loadGenerations()
			},
			error: err => {
				this.loading = false
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
	}

	loadGenerations() {
		if (!this.project) return
		this.gensLoading = true
		this.generationService
			.getAll(this.project.id, this.page, this.limit)
			.subscribe({
				next: gens => {
					this.generations = gens.generations || []
					this.totalCount = gens.totalCount || 0
					this.page = gens.page || 1
					this.loading = false
					this.gensLoading = false
				},
				error: () => {
					this.loading = false
					this.gensLoading = false
				},
			})
	}

	goToPage(p: number) {
		if (p < 1 || p > this.totalPages || p === this.page) return
		this.page = p
		this.loadGenerations()
	}

	startGeneration() {
		if (!this.project || this.modelsLoading || this.allModelsUnavailable) {
			return
		}

		if (!this.project) return
		this.generating = true
		this.projectService
			.startGeneration(this.project.id, this.selectedModel)
			.subscribe({
				next: gen => {
					this.generating = false
					this.router.navigate(['/generation', gen.id])
				},
				error: () => {
					this.generating = false
				},
			})
	}

	viewGeneration(gen: Generation) {
		if (gen.status === 'completed') {
			this.router.navigate(['/structure', gen.id])
		} else {
			this.router.navigate(['/generation', gen.id])
		}
	}

	toggleEdit() {
		this.editing = !this.editing
		if (this.editing) this.editForm = { ...this.project }
	}

	saveEdit() {
		if (!this.project) return
		this.saving = true
		this.projectService.edit(this.project.id, this.editForm).subscribe({
			next: () => {
				this.project = { ...this.project!, ...this.editForm }
				this.editing = false
				this.saving = false
			},
			error: () => {
				this.saving = false
			},
		})
	}

	deleteProject() {
		if (!this.project || !confirm('Delete this project?')) return
		this.projectService.delete(this.project.id).subscribe({
			next: () => this.router.navigate(['/projects']),
		})
	}

	formatDate(d: string) {
		return new Date(d).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
		})
	}

	statusColor(s: string) {
		return {
			'badge-pending': s === 'pending',
			'badge-process': s === 'process',
			'badge-completed': s === 'completed',
			'badge-failed': s === 'failed',
		}
	}

	min(a: number, b: number) {
		return Math.min(a, b)
	}
}
