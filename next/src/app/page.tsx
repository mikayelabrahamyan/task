import { getCreator, getCreators, getProducts, getSortedCreators } from "@/grpc";
import { GetCreatorsRequest, GetProductsRequest, GetSortedCreatorsRequest, SortOrder } from '@/grpc/marketplace_pb';
import Link from "next/link";

interface Props {
  searchParams: {
    order: string;
  }
}

export default async function Home({ searchParams: { order = "0" } }: Props) {
  // const request = new GetCreatorsRequest();
  // const creators = await getCreators(request);

  // const request1 = new GetProductsRequest();
  // const products = await getProducts(request1);

  const request3 = new GetSortedCreatorsRequest();
  request3.setLimit(3);
  request3.setOrder(parseInt(order) as SortOrder);
  const sortedCreators = await getSortedCreators(request3);
  

  return (
    <main className="flex min-h-screen flex-col items-start gap-[20px]">
      <nav className="flex gap-2">
        <Link href={'?order=0'}><button type="button">ASC</button></Link>
        <Link href={'?order=1'}><button type="button">DESC</button></Link>
      </nav>
      {
        sortedCreators.getCreatorsList().map(item => item.getEmail()).map(email => (<div key={email}>{email}</div>))
      }
    </main>
  );
}
