import Link, { type LinkProps } from 'next/link';

import { Button, type ButtonProps } from './button';

type LinkButtonProps = {
  href: LinkProps['href'];
  disabled?: boolean;
} & Omit<ButtonProps, 'asChild'> &
  Omit<LinkProps, 'href'>;

const LinkButton = ({ href, disabled, children, onClick, ...props }: LinkButtonProps) => {
  if (disabled) {
    return (
      <Button disabled {...props} >
        {children}
      </Button>
    );
  }

  return (
    <Button asChild {...props}>
      <Link
        href={href}
        onClick={(e) => {
          if (disabled) {
            e.preventDefault();
            return;
          }
          onClick?.(e);
        }}>
        {children}
      </Link>
    </Button>
  );
};

LinkButton.displayName = 'LinkButton';

export { LinkButton };
