import { Component } from '@angular/core'
import { RouterOutlet } from '@angular/router'
import { NavbarComponent } from './components/navbar/navbar.component'

@Component({
	selector: 'app-root',
	standalone: true,
	imports: [RouterOutlet, NavbarComponent],
	template: `
		<app-navbar />
		<main class="main-content">
			<router-outlet />
		</main>
	`,
	styles: [
		`
			:host {
				display: block;
				min-height: 100vh;
			}
			.main-content {
				padding-top: 64px;
			}
		`,
	],
})
export class AppComponent {}
