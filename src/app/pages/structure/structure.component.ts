import { CommonModule } from '@angular/common'
import { HttpClient } from '@angular/common/http'
import { Component, inject, OnInit } from '@angular/core'
import { ActivatedRoute, Router, RouterLink } from '@angular/router'
import { Template } from '../../models'
import { AuthService } from '../../services/auth.service'
import { GenerationService } from '../../services/generation.service'

interface TreeNode {
	name: string
	type: 'file' | 'dir'
	children?: TreeNode[]
	isFile?: boolean
}

@Component({
	selector: 'app-structure',
	standalone: true,
	imports: [CommonModule, RouterLink],
	templateUrl: './structure.component.html',
	styleUrls: ['./structure.component.scss'],
})
export class StructureComponent implements OnInit {
	route = inject(ActivatedRoute)
	router = inject(Router)
	http = inject(HttpClient)
	generationService = inject(GenerationService)
	authService = inject(AuthService)

	templates: Template[] = []
	activeTemplate: Template | null = null
	activeTab: 'simple' | 'medium' | 'enterprise' = 'simple'
	loading = true
	tree: TreeNode[] = []
	downloading = false
	expandedNodes = new Set<string>()

	typeLabels: Record<string, { label: string; desc: string; icon: string }> = {
		simple: {
			label: 'Simple',
			desc: 'Clean & straightforward structure for small projects',
			icon: '◇',
		},
		medium: {
			label: 'Medium',
			desc: 'Organized structure for growing applications',
			icon: '◈',
		},
		enterprise: {
			label: 'Enterprise',
			desc: 'Advanced modular structure for large-scale systems',
			icon: '◆',
		},
	}

	ngOnInit() {
		this.authService.getMe().subscribe({
			error: err => {
				if (typeof err === 'object' && 'status' in err && err.status === 401) {
					this.router.navigate(['/auth'])
				}
			},
		})
		const genId = this.route.snapshot.paramMap.get('genId')!
		this.generationService.getTemplates(genId).subscribe({
			next: tpls => {
				this.templates = tpls || []
				if (this.templates.length > 0) {
					const first =
						this.templates.find(t => t.type === 'simple') || this.templates[0]
					this.selectTemplate(first)
				}
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

	selectTemplate(t: Template) {
		this.activeTemplate = t
		this.activeTab = t.type as any
		this.tree = this.buildTree(t.content.files, t.content.directories)
		this.expandedNodes.clear()
		// auto-expand first level
		this.tree.forEach(n => this.expandedNodes.add(n.name))
	}

	buildTree(files: string[], dirs: string[]): TreeNode[] {
		const root: Record<string, any> = {}

		const ensureDir = (parts: string[], node: any, path: string) => {
			if (!parts.length) return
			const part = parts[0]
			const currentPath = path ? `${path}/${part}` : part
			if (!node[part])
				node[part] = { __type: 'dir', __path: currentPath, __children: {} }
			ensureDir(parts.slice(1), node[part].__children, currentPath)
		}

		const addFile = (parts: string[], node: any, path: string) => {
			if (parts.length === 1) {
				const currentPath = path ? `${path}/${parts[0]}` : parts[0]
				node[parts[0]] = { __type: 'file', __path: currentPath }
				return
			}
			const part = parts[0]
			const currentPath = path ? `${path}/${part}` : part
			if (!node[part])
				node[part] = { __type: 'dir', __path: currentPath, __children: {} }
			addFile(parts.slice(1), node[part].__children, currentPath)
		}

		dirs.forEach(d => ensureDir(d.split('/'), root, ''))
		files.forEach(f => addFile(f.split('/'), root, ''))

		const toNodes = (node: any): TreeNode[] => {
			return Object.entries(node)
				.filter(([k]) => !k.startsWith('__'))
				.sort(([, a]: any, [, b]: any) => {
					if (a.__type !== b.__type) return a.__type === 'dir' ? -1 : 1
					return 0
				})
				.map(([name, val]: any) => ({
					name,
					type: val.__type,
					path: val.__path,
					children: val.__type === 'dir' ? toNodes(val.__children) : undefined,
				}))
		}

		return toNodes(root)
	}

	toggleNode(path: string) {
		if (this.expandedNodes.has(path)) this.expandedNodes.delete(path)
		else this.expandedNodes.add(path)
	}

	isExpanded(path: string) {
		return this.expandedNodes.has(path)
	}

	download() {
		if (!this.activeTemplate) return
		this.downloading = true
		const url = this.generationService.getDownloadUrl(
			this.activeTemplate.generationId
		)
		// Try to trigger browser download with template type as param
		const link = document.createElement('a')
		link.href = `${url}?type=${this.activeTemplate.type}`
		link.setAttribute('download', '')
		const token = localStorage.getItem('token')
		// Use fetch for auth header
		fetch(link.href, { headers: { Authorization: `Bearer ${token}` } })
			.then(r => r.blob())
			.then(blob => {
				const objUrl = URL.createObjectURL(blob)
				link.href = objUrl
				link.download = `project-structure-${this.activeTemplate!.type}.zip`
				document.body.appendChild(link)
				link.click()
				document.body.removeChild(link)
				URL.revokeObjectURL(objUrl)
				this.downloading = false
			})
			.catch(() => {
				this.downloading = false
			})
	}

	getFileIcon(name: string): string {
		if (name.endsWith('.go')) return '🔵'
		if (name.endsWith('.ts') || name.endsWith('.js')) return '🟡'
		if (name.endsWith('.sql')) return '🟣'
		if (name.endsWith('.md')) return '📄'
		if (name.endsWith('.yaml') || name.endsWith('.yml')) return '🟠'
		if (name.endsWith('.json')) return '🟤'
		if (name === 'Dockerfile') return '🐳'
		if (name === 'Makefile') return '⚙️'
		if (name.startsWith('.')) return '⚫'
		return '📄'
	}

	get totalFiles(): number {
		return this.activeTemplate?.content.files.length || 0
	}
	get totalDirs(): number {
		return this.activeTemplate?.content.directories.length || 0
	}
}
