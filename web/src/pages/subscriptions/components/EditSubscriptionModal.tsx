import { useState, useEffect } from 'react'
import { X, Loader2, Github, Rss, Database, Globe, Webhook, Shield, Bug, AlertTriangle } from 'lucide-react'
import { subscriptionService } from '@/services/subscription'
import toast from 'react-hot-toast'
import type { Subscription } from '@/types'

interface EditSubscriptionModalProps {
  isOpen: boolean
  onClose: () => void
  onSuccess?: () => void
  subscription: Subscription | null
}

const sourceTypeIcons: Record<string, typeof Github> = {
  github_repo: Github,
  rss: Rss,
  nvd: Database,
  cve: Bug,
  threat_intel: Shield,
  vendor_advisory: Globe,
  attack_activity: AlertTriangle,
  webhook: Webhook,
}

const sourceTypeLabels: Record<string, string> = {
  github_repo: 'GitHub 仓库',
  rss: 'RSS 订阅',
  nvd: 'NVD 漏洞库',
  cve: 'CVE 数据',
  threat_intel: '威胁情报',
  vendor_advisory: '厂商公告',
  attack_activity: '攻击活动',
  webhook: 'Webhook',
}

// 将分钟转换为 cron 表达式
function minutesToCron(minutes: number): string {
  if (minutes < 60) {
    return `*/${minutes} * * * *`
  }
  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    return `0 */${hours} * * *`
  }
  return '0 0 * * *'
}

// 将 cron 表达式转换为分钟
function cronToMinutes(cron: string): number {
  const parts = cron.split(' ')
  if (parts[0].startsWith('*/')) {
    return parseInt(parts[0].substring(2))
  }
  if (parts[1].startsWith('*/')) {
    return parseInt(parts[1].substring(2)) * 60
  }
  return 60
}

export default function EditSubscriptionModal({
  isOpen,
  onClose,
  onSuccess,
  subscription,
}: EditSubscriptionModalProps) {
  const [formData, setFormData] = useState({
    name: '',
    source_url: '',
    fetch_interval: 60,
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (subscription) {
      setFormData({
        name: subscription.name,
        source_url: subscription.source_url,
        fetch_interval: cronToMinutes(subscription.cron_expr),
      })
    }
  }, [subscription])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!subscription) return

    setIsSubmitting(true)
    try {
      await subscriptionService.update(subscription.id, {
        name: formData.name,
        source_url: formData.source_url,
        cron_expr: minutesToCron(formData.fetch_interval),
      })
      toast.success('订阅更新成功')
      onSuccess?.()
      onClose()
    } catch (error) {
      toast.error('更新订阅失败')
      console.error(error)
    } finally {
      setIsSubmitting(false)
    }
  }

  if (!isOpen || !subscription) return null

  const SourceIcon = sourceTypeIcons[subscription.source_type] || Database

  return (
    <>
      {/* Backdrop */}
      <div className="modal-overlay" onClick={onClose}>
        {/* Modal */}
        <div className="modal w-full max-w-lg" onClick={(e) => e.stopPropagation()}>
          {/* Header */}
          <div className="modal-header">
            <div>
              <h2 className="modal-title">编辑订阅</h2>
              <p className="text-sm text-gray-500 mt-0.5">
                修改订阅配置信息
              </p>
            </div>
            <button
              onClick={onClose}
              className="p-1.5 rounded text-gray-400 hover:text-gray-200 hover:bg-gray-800 transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="modal-body max-h-[60vh] space-y-4">
              {/* Source Type Badge */}
              <div className="flex items-center gap-3 p-3 rounded-lg bg-gray-800">
                <SourceIcon className="w-5 h-5 text-primary-400" />
                <span className="text-sm text-gray-200">
                  {sourceTypeLabels[subscription.source_type] || subscription.source_type}
                </span>
                <span className="ml-auto text-xs text-gray-500">数据源类型</span>
              </div>

              {/* Name */}
              <div className="form-item">
                <label className="label">订阅名称 <span className="text-danger-500">*</span></label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="输入订阅名称"
                  className="input"
                  required
                />
              </div>

              {/* URL */}
              {subscription.source_type !== 'webhook' && subscription.source_type !== 'cve' && (
                <div className="form-item">
                  <label className="label">订阅地址</label>
                  <input
                    type="text"
                    value={formData.source_url}
                    onChange={(e) => setFormData({ ...formData, source_url: e.target.value })}
                    placeholder="https://example.com/feed"
                    className="input"
                  />
                </div>
              )}

              {/* Fetch Interval */}
              <div className="form-item">
                <label className="label">抓取频率 <span className="text-danger-500">*</span></label>
                <select
                  value={formData.fetch_interval}
                  onChange={(e) => setFormData({ ...formData, fetch_interval: Number(e.target.value) })}
                  className="input"
                >
                  <option value={30}>每 30 分钟</option>
                  <option value={60}>每 1 小时</option>
                  <option value={180}>每 3 小时</option>
                  <option value={360}>每 6 小时</option>
                  <option value={720}>每 12 小时</option>
                  <option value={1440}>每 24 小时</option>
                </select>
                <p className="text-xs text-gray-500 mt-1.5">
                  Cron 表达式: {minutesToCron(formData.fetch_interval)}
                </p>
              </div>
            </div>

            {/* Footer */}
            <div className="modal-footer">
              <button
                type="button"
                onClick={onClose}
                className="btn-default"
                disabled={isSubmitting}
              >
                取消
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="btn-primary"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin" />
                    保存中...
                  </>
                ) : (
                  '保存更改'
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  )
}
