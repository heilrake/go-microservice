import * as React from 'react'
import clsx from 'clsx'

type InputProps = React.ComponentProps<'input'> & {
  label?: string
  error?: string
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className, type = "text", ...props }, ref) => {
    return (
      <div className="flex flex-col w-full">
        {label && (
          <label className="mb-1 text-sm font-medium text-gray-700">
            {label}
          </label>
        )}
        <input
          type={type}
          ref={ref}
          className={clsx(
            'flex h-10 w-full rounded-md border px-3 py-2 text-base ring-offset-background placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
            'border-input bg-background text-foreground focus:border-primary focus:ring-1 focus:ring-primary',
            error && 'border-red-500 focus:ring-red-500',
            className
          )}
          {...props}
        />
        {error && (
          <p className="mt-1 text-sm text-red-600">{error}</p>
        )}
      </div>
    )
  }
)

Input.displayName = 'Input'

export { Input }
