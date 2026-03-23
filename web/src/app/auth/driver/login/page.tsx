import { redirect } from 'next/navigation';

export default function DriverLoginPage() {
  redirect('/auth?role=driver');
}
