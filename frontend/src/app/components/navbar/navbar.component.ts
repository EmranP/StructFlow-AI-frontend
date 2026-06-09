import { CommonModule } from '@angular/common'
import { Component, inject } from '@angular/core'
import { Router, RouterLink, RouterLinkActive } from '@angular/router'
import { AuthService } from '../../services/auth.service'

@Component({
	selector: 'app-navbar',
	standalone: true,
	imports: [CommonModule, RouterLink, RouterLinkActive],
	templateUrl: './navbar.component.html',
	styleUrls: ['./navbar.component.scss'],
})
export class NavbarComponent {
	auth = inject(AuthService)
	router = inject(Router)
	loggingOut = false

	logout() {
		this.loggingOut = true
		this.auth.logout().subscribe({
			next: () => {
				this.loggingOut = false
				this.router.navigate(['/auth'])
			},
			error: () => {
				this.loggingOut = false
				this.router.navigate(['/auth'])
			},
		})
	}
}
