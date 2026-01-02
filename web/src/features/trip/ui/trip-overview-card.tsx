import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card"
import { ScrollArea } from "@/shared/ui/scroll-area"

type TripOverviewCardProps = {
  title: string
  description: string
  children?: React.ReactNode
}

export const TripOverviewCard = ({ title, description, children }: TripOverviewCardProps) => {
  return (
    <Card className="w-full md:max-w-[500px] z-[9999] flex-[0.3]">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        <ScrollArea>
          {children}
        </ScrollArea>
      </CardContent>
    </Card>
  )
}

