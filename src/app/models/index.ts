export interface User {
	id: string
	email: string
	createdAt: string
}

export interface MessageResponse {
	message: string
}

export interface AuthResponse {
	accessToken: string
}

export interface GetID {
	id: string
}

export interface ProjectItem {
	id: string
	title: string
	projectType: string
	stack: string
	architecture: string
	features: string
	additionalInfo: string
	createdAt: string
	updatedAt: string
}

export interface Projects {
	limit: number
	page: number
	projects: ProjectItem[]
	totalCount: number
}

export interface CreateProjectDto {
	title: string
	projectType: string
	stack: string
	architecture: string
	features: string
	additionalInfo: string
}

export interface Generation {
	id: string
	projectId: string
	status: 'pending' | 'process' | 'completed' | 'failed'
	errorMessage?: string
	createdAt: string
	updatedAt: string
}

export interface GenerationResponse {
	limit: number
	page: number
	totalCount: number
	generations: Generation[]
}

export interface TemplateContent {
	files: string[]
	directories: string[]
}

export interface Template {
	id: string
	generationId: string
	type: 'simple' | 'medium' | 'enterprise'
	content: TemplateContent
	createdAt: string
}

export interface ModelResponse {
	id: string
	name: string
	available: boolean
}

export interface VerifyEmailDto {
	email: string
	code: string
}

export interface RefreshResponse {
	token: string
}
