import { CommonModule } from '@angular/common'
import { Component, inject, OnInit } from '@angular/core'
import { FormsModule } from '@angular/forms'
import { Router, RouterLink } from '@angular/router'
import { CreateProjectDto, ProjectItem } from '../../models'
import { AuthService } from '../../services/auth.service'
import { ProjectService } from '../../services/project.service'

@Component({
	selector: 'app-projects',
	standalone: true,
	imports: [CommonModule, RouterLink, FormsModule],
	templateUrl: './projects.component.html',
	styleUrls: ['./projects.component.scss', '../../app.component.scss'],
})
export class ProjectsComponent implements OnInit {
	router = inject(Router)
	projectService = inject(ProjectService)
	authService = inject(AuthService)

	projects: ProjectItem[] = []
	loading = true
	showModal = false
	creating = false
	error = ''

	// Pagination
	page = 1
	limit = 10
	totalCount = 0

	form: CreateProjectDto = {
		title: '',
		projectType: '',
		stack: '',
		architecture: '',
		features: '',
		additionalInfo: '',
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
		this.loadProjects()
	}

	loadProjects() {
		this.loading = true
		this.projectService.getAll(this.page, this.limit).subscribe({
			next: data => {
				this.projects = data.projects || []
				this.totalCount = data.totalCount || 0
				this.page = data.page || 1
				this.loading = false
			},
			error: err => {
				this.loading = false
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
	}

	goToPage(p: number) {
		if (p < 1 || p > this.totalPages || p === this.page) return
		this.page = p
		this.loadProjects()
	}

	openModal() {
		this.showModal = true
		this.resetForm()
	}
	closeModal() {
		this.showModal = false
	}

	resetForm() {
		this.form = {
			title: '',
			projectType: '',
			stack: '',
			architecture: '',
			features: '',
			additionalInfo: '',
		}
		this.error = ''
	}

	createProject() {
		if (!this.form.title.trim()) {
			this.error = 'Title is required'
			return
		}
		this.creating = true
		this.error = ''
		this.projectService.create(this.form).subscribe({
			next: project => {
				this.creating = false
				this.showModal = false
				this.router.navigate(['/projects', project.id])
			},
			error: e => {
				this.creating = false
				this.error = e?.error?.message || 'Failed to create project'
			},
		})
	}

	deleteProject(id: string, e: Event) {
		e.stopPropagation()
		e.preventDefault()
		if (!confirm('Delete this project?')) return
		this.projectService.delete(id).subscribe({
			next: () => {
				this.projects = this.projects.filter(p => p.id !== id)
				this.totalCount--
				if (this.projects.length === 0 && this.page > 1) {
					this.page--
					this.loadProjects()
				}
			},
		})
	}

	formatDate(d: string) {
		return new Date(d).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric',
		})
	}

	min(a: number, b: number) {
		return Math.min(a, b)
	}
}
