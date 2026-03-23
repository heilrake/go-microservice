import { redirect } from 'next/navigation';

export default function RiderLoginPage() {
  redirect('/auth?role=rider');
}
