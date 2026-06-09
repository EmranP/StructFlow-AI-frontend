import { Routes } from '@angular/router'
import { authGuard } from './guards/auth.guard'

export const routes: Routes = [
	{ path: '', redirectTo: '/projects', pathMatch: 'full' },
	{
		path: 'auth',
		loadComponent: () =>
			import('./pages/auth/auth.component').then(m => m.AuthComponent),
	},
	{
		path: 'verify-email',
		loadComponent: () =>
			import('./pages/verify-email/verify-email.component').then(
				m => m.VerifyEmailComponent
			),
	},
	{
		path: 'projects',
		loadComponent: () =>
			import('./pages/projects/projects.component').then(
				m => m.ProjectsComponent
			),
		canActivate: [authGuard],
	},
	{
		path: 'projects/:id',
		loadComponent: () =>
			import('./pages/project-detail/project-detail.component').then(
				m => m.ProjectDetailComponent
			),
		canActivate: [authGuard],
	},
	{
		path: 'generation/:id',
		loadComponent: () =>
			import('./pages/generation/generation.component').then(
				m => m.GenerationComponent
			),
		canActivate: [authGuard],
	},
	{
		path: 'structure/:genId',
		loadComponent: () =>
			import('./pages/structure/structure.component').then(
				m => m.StructureComponent
			),
		canActivate: [authGuard],
	},
	{ path: '**', redirectTo: '/projects' },
]
