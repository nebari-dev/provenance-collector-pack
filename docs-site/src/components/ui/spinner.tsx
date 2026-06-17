// Landed via `shadcn add @nebari/spinner`. Source: nebari-design.
// Do not edit by hand — re-run the shadcn command to refresh.
import { cva, type VariantProps } from 'class-variance-authority';
import { LoaderCircle, type LucideProps } from 'lucide-react';
import { cn } from '@/lib/utils';

const spinnerVariants = cva('animate-spin', {
  variants: {
    size: {
      xs: 'size-3.5',
      sm: 'size-4',
      default: '',
      lg: 'size-6',
      xl: 'size-8',
    },
  },
  defaultVariants: {
    size: 'default',
  },
});

type SpinnerProps = Omit<LucideProps, 'size'> &
  VariantProps<typeof spinnerVariants> & {
    label?: string;
  };

function Spinner({
  className,
  size,
  label = 'Loading',
  ...props
}: SpinnerProps) {
  return (
    <LoaderCircle
      aria-label={label}
      className={cn(spinnerVariants({ size }), className)}
      data-size={size ?? 'default'}
      data-slot="spinner"
      role="status"
      {...props}
    />
  );
}

export type { SpinnerProps };
export { Spinner, spinnerVariants };
